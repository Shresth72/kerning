import asyncio
from pathlib import Path

from utils import goto, skip_bytes, read_uint16, read_uint32, read_tag
from glyph import (
    ReadSimpleGlyph,
    GlyphDrawTest,
    GetAllGlyphLocations,
    DrawTextGlyphs,
    GetUnicodeToGlyphIndexMappings,
)


def RenderText(f, text, glyph_index_map, all_glyph_locs):
    glyphs_to_draw = []
    for c in text:
        if c == " ":
            glyphs_to_draw.append((None, True))
            continue

        unicode = ord(c)
        glyph_index = glyph_index_map.get(unicode, None)

        if glyph_index is not None and glyph_index < len(all_glyph_locs):
            glyph_loc = all_glyph_locs[glyph_index]
            goto(f, glyph_loc)
            glyph_data = ReadSimpleGlyph(f)
            glyphs_to_draw.append((glyph_data, False))

    # TODO: Do this heavy lifting in another rust library
    DrawTextGlyphs(glyphs_to_draw)


def RenderAllGlyphs(f, all_glyph_locs):
    for i, glyph_loc in enumerate(all_glyph_locs):
        goto(f, glyph_loc)
        glyph_data = ReadSimpleGlyph(f)
        GlyphDrawTest(glyph_data, f"glyfs/glyph-{i}.png")
        if i == 20:
            break
    print("Saved all glyphs to glyfs")


async def ParseFont(fontPath: str) -> None:
    Path("glyfs").mkdir(exist_ok=True)

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

        unicode_to_glyph_index_map = GetUnicodeToGlyphIndexMappings(
            f, table_location_map
        )

        text = "Shrestha"
        RenderText(f, text, unicode_to_glyph_index_map, all_glyph_locs)
        # RenderAllGlyphs(f, all_glyph_locs)


if __name__ == "__main__":
    asyncio.run(ParseFont("../assets/Meditative.ttf"))
