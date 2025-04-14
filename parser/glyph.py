import matplotlib.pyplot as plt
from typing import BinaryIO, List

from bezier import Point, DrawBezier
from utils import (
    get_location,
    goto,
    read_uint8,
    read_int16,
    read_uint16,
    read_uint32,
    skip_bytes,
    flag_bit_is_set,
)

MAX_UINT = 0xFFFFFFFF


class GlyphData:
    def __init__(self, coordsX: list, coordsY: list, contour_end_indices: list):
        self.points = [Point(x, y) for x, y in zip(coordsX, coordsY)]
        self.coordsX = coordsX
        self.coordsY = coordsY
        self.contour_end_indices = contour_end_indices

    def __str__(self):
        return (
            f"GlyphData(\n"
            f"    points={self.points},\n"
            f"    contour_end_indices={self.contour_end_indices}\n)"
        )


def DrawTextGlyphs(
    glyphs: list, spacing: int = 700, word_spacing: int = 500, filename="glyfs/text.png"
):
    fig, ax = plt.subplots()
    ax.set_aspect("equal")
    ax.set_title("Text Rendering")

    x_cursor = 0

    for glyph, is_space in glyphs:
        if is_space:
            x_cursor += word_spacing
            continue

        resolution = 100
        contour_start_index = 0
        for contour_end_index in glyph.contour_end_indices:
            num_points_in_contour = contour_end_index - contour_start_index + 1

            points = [
                Point(
                    float(glyph.coordsX[contour_start_index + i]) + x_cursor,
                    float(glyph.coordsY[contour_start_index + i]),
                )
                for i in range(num_points_in_contour)
            ]

            for i in range(0, len(points), 2):
                p0 = points[i]
                p1 = points[(i + 1) % len(points)]
                p2 = points[(i + 2) % len(points)]
                DrawBezier(p0, p1, p2, resolution, ax)

            contour_start_index = contour_end_index + 1

        x_cursor += spacing

    plt.savefig(filename)
    print(f"Saved rendered text to {filename}")
    plt.close(fig)


def GlyphDrawTest(glyph: GlyphData, filename="glyfs/glyph_output.png"):
    fig, ax = plt.subplots()
    ax.set_aspect("equal")
    ax.set_title("Glyph Drawing")

    resolution = 100
    contour_start_index = 0
    for contour_end_index in glyph.contour_end_indices:
        num_points_in_contour = contour_end_index - contour_start_index + 1

        points = [
            Point(
                float(glyph.coordsX[contour_start_index + i]),
                float(glyph.coordsY[contour_start_index + i]),
            )
            for i in range(num_points_in_contour)
        ]

        for i in range(0, len(points), 2):
            p0 = points[i]
            p1 = points[(i + 1) % len(points)]
            p2 = points[(i + 2) % len(points)]
            DrawBezier(p0, p1, p2, resolution, ax)

        contour_start_index = contour_end_index + 1

    # for point in glyph.points:
    #     DrawPoint(point, ax)

    plt.savefig(filename)
    print(f"Saved glyph to {filename}")
    plt.close(fig)


def GetAllGlyphLocations(f: BinaryIO, lookup: dict):
    goto(f, lookup["maxp"] + 4)
    num_glyphs = read_uint16(f)

    goto(f, lookup["head"])
    skip_bytes(f, 50)

    is_two_byte_entry = True if read_uint16(f) == 0 else False

    location_table_start = lookup["loca"]
    glyph_table_start = lookup["glyf"]

    all_glyph_locs = [0] * num_glyphs
    for glyph_index in range(num_glyphs):
        goto(f, location_table_start + glyph_index * (2 if is_two_byte_entry else 4))
        glyph_data_offset = read_uint16(f) * 2 if is_two_byte_entry else read_uint32(f)
        all_glyph_locs[glyph_index] = glyph_table_start + glyph_data_offset

    return all_glyph_locs


def ReadCoordinates(f: BinaryIO, all_flags: list, readingX: bool) -> list:
    offsetSizeFlagBit = 1 if readingX else 2
    offsetSignOrSkipBit = 4 if readingX else 5
    coordinates = [0] * len(all_flags)

    for i in range(len(coordinates)):
        coordinates[i] = coordinates[max(0, i - 1)]
        flag = all_flags[i]
        # on_curve = flag_bit_is_set(flag, 0)

        if flag_bit_is_set(flag, offsetSizeFlagBit):
            offset = read_uint8(f)
            sign = 1 if flag_bit_is_set(flag, offsetSignOrSkipBit) else -1
            coordinates[i] += offset * sign
        elif not flag_bit_is_set(flag, offsetSignOrSkipBit):
            coordinates[i] += read_int16(f)

    return coordinates


def ReadSimpleGlyph(f: BinaryIO) -> GlyphData:
    num_contours = read_uint16(f)
    contour_end_indices = [0] * num_contours
    skip_bytes(f, 8)  # Skip bound size

    for i in range(num_contours):
        contour_end_indices[i] = read_uint16(f)

    # Read flags
    num_points = contour_end_indices[-1] + 1
    all_flags = [0] * num_points
    skip_bytes(f, read_uint16(f))  # Skip instructuions

    i = 0
    while i < num_points:
        flag = read_uint8(f)
        all_flags[i] = flag

        if flag_bit_is_set(flag, 3):
            for r in range(read_uint8(f)):
                i += 1
                if i < num_points:
                    all_flags[i] = flag
        i += 1

    coordsX = ReadCoordinates(f, all_flags, readingX=True)
    coordsY = ReadCoordinates(f, all_flags, readingX=False)
    return GlyphData(coordsX, coordsY, contour_end_indices)


def GetUnicodeToGlyphIndexMappings(f: BinaryIO, lookUp: dict) -> dict:
    """
    params:
        - Binary File reader
        - Lookup table for {tag: offset}

    returns:
        Dict of {unicode: glyph_index}
    """
    mappings = []
    cmap_offset = lookUp["cmap"]
    goto(f, cmap_offset)

    version = read_uint16(f)
    num_subtables = read_uint16(f)

    cmap_subtable_offset = MAX_UINT

    for i in range(num_subtables):
        platform_id = read_uint16(f)
        platform_specific_id = read_uint16(f)
        offset = read_uint32(f)

        # platform_id of 0 means Unicode, thus platform_specific_id can be interpreted as Unicode version
        if platform_id == 0:
            unicode_version_info = platform_specific_id

            # Unicode>=2.0 semantics (non-BMP chars allowed)
            if unicode_version_info == 4:
                cmap_subtable_offset = offset
            # Unicode>=2.0 semantics (BMP only)
            if unicode_version_info == 3 and cmap_subtable_offset == MAX_UINT:
                cmap_subtable_offset = offset

    if cmap_subtable_offset == 0:
        raise Exception("Font does not contain supported character map type (TODO)")

    goto(f, cmap_offset + cmap_subtable_offset)
    format = read_uint16(f)

    if format != 12 and format != 4:
        raise Exception(f"Font cmap format not supported (TODO): {format}")
    elif format == 12:
        read_uint16(f)  # reserved
        subtable_byte_length_including_header = read_uint32(f)
        language_code = read_uint32(f)
        num_groups = read_uint32(f)

        for i in range(num_groups):
            start_char_code = read_uint32(f)
            end_char_code = read_uint32(f)
            start_glyph_index = read_uint32(f)

            num_chars = end_char_code - start_char_code + 1
            for char_code_offset in range(num_chars):
                char_code = start_char_code + char_code_offset
                glyph_index = start_glyph_index + char_code_offset
                mappings.append((char_code, glyph_index))
    elif format == 4:
        skip_bytes(f, 4)  # Skip length, language_code

        # Number of continuous segments of character codes
        seg_count_2x = read_uint16(f)
        seg_count = int(seg_count_2x / 2)
        skip_bytes(f, 6)  # Skip search_range, entry_selector, range_shift

        end_codes = [read_uint16(f) for _ in range(seg_count)]
        skip_bytes(f, 2)  # Reserved pad
        start_codes = [read_uint16(f) for _ in range(seg_count)]
        id_deltas = [read_uint16(f) for _ in range(seg_count)]
        id_range_offsets = [
            (get_location(f), read_uint16(f)) for _ in range(seg_count)
        ]  # (readLoc, offset)

        for i in range(len(start_codes)):
            end_code = end_codes[i]
            curr_code = start_codes[i]
            while curr_code <= end_code:
                glyph_index = 0
                if id_range_offsets[i][1] == 0:  # offset
                    glyph_index = (curr_code + id_deltas[i]) % 65536
                else:
                    range_offset_location = (
                        id_range_offsets[i][0] + id_range_offsets[i][1]
                    )
                    glyph_index_array_location = (
                        2 * (curr_code - start_codes[i]) + range_offset_location
                    )

                    reader_location_old = get_location(f)
                    goto(f, glyph_index_array_location)
                    glyph_index_offset = read_uint16(f)
                    goto(f, reader_location_old)

                    if glyph_index_offset != 0:
                        glyph_index = (glyph_index_offset + id_deltas[i]) % 65536

                mappings.append((curr_code, glyph_index))
                curr_code += 1

    return dict(mappings)
