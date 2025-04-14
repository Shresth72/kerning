from pathlib import Path

from utils import goto, skip_bytes, read_uint16, read_uint32, read_tag
from glyph import ReadSimpleGlyph, GlyphDrawTest, GetAllGlyphLocations


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

        all_glyph_locs = GetAllGlyphLocations(f, table_location_map)

        Path("glyfs").mkdir(exist_ok=True)

        for i, glyph_loc in enumerate(all_glyph_locs):
            goto(f, glyph_loc)
            glyph_data = ReadSimpleGlyph(f)
            GlyphDrawTest(glyph_data, f"glyfs/glyph-{i}.png")

        print("Saved all glyphs to glyfs")


if __name__ == "__main__":
    ParseFont("../assets/Meditative.ttf")
