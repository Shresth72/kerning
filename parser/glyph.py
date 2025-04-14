from typing import BinaryIO
from utils import read_uint8, read_uint16, skip_bytes, flag_bit_is_set


class Point:
    def __init__(self, x: int, y: int):
        self.x = x
        self.y = y

    def __repr__(self):
        return f"({self.x}, {self.y})"


class GlyphData:
    def __init__(self, coordsX: list, coordsY: list, contour_end_indices: list):
        self.points = [Point(x, y) for x, y in zip(coordsX, coordsY)]
        self.contour_end_indices = contour_end_indices

    def __str__(self):
        return (
            f"GlyphData(\n"
            f"    points={self.points},\n"
            f"    contour_end_indices={self.contour_end_indices}\n)"
        )


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
            coordinates[i] += read_uint16(f)

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
