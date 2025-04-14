from utils import goto, skip_bytes, read_uint16, read_uint32, read_tag
from glyph import ReadSimpleGlyph


def ParseFont(fontPath: str) -> None:
    with open(fontPath, "rb") as f:
        skip_bytes(f, 4)
        num_tables = read_uint16(f)
        skip_bytes(f, 6)
        print(f"NumTables: {num_tables}")

        table_location_map = dict()

        for i in range(num_tables):
            tag = read_tag(f)
            skip_bytes(f, 4)  # checksum
            offset = read_uint32(f)
            skip_bytes(f, 4)  # length

            table_location_map[tag] = offset

        if "glyf" in table_location_map:
            glyf_offset = table_location_map["glyf"]
            goto(f, glyf_offset)

            glyf0 = ReadSimpleGlyph(f)
            print(f"Glyf0: \n{glyf0}")


if __name__ == "__main__":
    ParseFont("../assets/JetBrainsMono-Bold.ttf")
