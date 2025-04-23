#include "utils.h"

void GetAllGlyphLocations(FILE *f, GlyphTableMap *glyph_map) {
  goto_loc(f, get_table_offset(glyph_map, "maxp") + 4);
  __uint16_t num_glyphs = read_uint16(f);
  printf("NumGlyphs: %d\n", num_glyphs);

  goto_loc(f, get_table_offset(glyph_map, "head"));
  skip_bytes(f, 50);

  bool is_two_byte_entry = (read_uint16(f) == 0);
  printf("is_two_byte_entry: %s\n", is_two_byte_entry ? "true" : "false");

  __uint32_t loc_table_start = get_table_offset(glyph_map, "loca");
  __uint32_t glyph_table_start = get_table_offset(glyph_map, "glyf");
}
