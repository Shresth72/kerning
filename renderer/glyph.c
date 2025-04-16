#include "utils.h"

void GetAllGlyphLocations(FILE *f, TableLoc *lookup, int table_count) {
  goto_loc(f, get_table_offset(lookup, table_count, "maxp") + 4);
  __uint16_t num_glyphs = read_uint16(f);

  goto_loc(f, get_table_offset(lookup, table_count, "head"));
  skip_bytes(f, 50);
}
