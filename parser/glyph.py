import matplotlib.pyplot as plt

from typing import BinaryIO
from utils import (
    goto,
    read_uint8,
    read_int16,
    read_uint16,
    read_uint32,
    skip_bytes,
    flag_bit_is_set,
)


class Point:
    def __init__(self, x: int, y: int):
        self.x = x
        self.y = y

    def __repr__(self):
        return f"({self.x}, {self.y})"


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


def DrawLine(p1: Point, p2: Point, ax):
    ax.plot([p1.x, p2.x], [p1.y, p2.y], "k-")


def DrawPoint(p: Point, ax):
    ax.plot(p.x, p.y, "ro", markersize=3)


def GlyphDrawTest(glyph: GlyphData, filename="glyfs/glyph_output.png"):
    fig, ax = plt.subplots()
    ax.set_aspect("equal")
    ax.set_title("Glyph Drawing")

    contour_start_index = 0
    for contour_end_index in glyph.contour_end_indices:
        num_points_in_contour = contour_end_index - contour_start_index + 1
        # points = glyph.points[contour_start_index:num_points_in_contour]

        points = [
            Point(
                float(glyph.coordsX[contour_start_index + i]),
                float(glyph.coordsY[contour_start_index + i]),
            )
            for i in range(num_points_in_contour)
        ]

        for i in range(len(points)):
            DrawLine(points[i], points[(i + 1) % len(points)], ax)
        contour_start_index = contour_end_index + 1

    for point in glyph.points:
        DrawPoint(point, ax)

    ax.invert_yaxis()
    plt.savefig(filename)
    # print(f"Saved glyph to {filename}")
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
