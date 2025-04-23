#include "utils.h"

long get_loc(FILE *file) { return ftell(file); }

void goto_loc(FILE *file, long position) { fseek(file, position, SEEK_SET); }

void skip_bytes(FILE *file, long num) { fseek(file, num, SEEK_CUR); }

__uint8_t read_uint8(FILE *file) {
  __uint8_t val;
  fread(&val, 1, 1, file);
  return val;
}

__int16_t read_int16(FILE *file) {
  __int8_t bytes[2];
  fread(&bytes, 1, 2, file);
  return (__int16_t)((bytes[0] << 8) | bytes[1]);
}

__uint16_t read_uint16(FILE *file) {
  __uint8_t bytes[2];
  fread(&bytes, 1, 2, file);
  return (__uint16_t)((bytes[0] << 8) | bytes[1]);
}

__uint32_t read_uint32(FILE *file) {
  __uint8_t bytes[4];
  fread(&bytes, 1, 4, file);
  return (__uint32_t)((bytes[0] << 24) | (bytes[1] << 16) | (bytes[2] << 8) |
                      bytes[3]);
}

void read_tag(FILE *file, char *tag) {
  fread(tag, 1, 4, file);
  tag[4] = '\0';
}

bool flag_bit_is_set(__uint32_t flag, int bit_index) {
  return ((flag >> bit_index) & 1) == 1;
}

__uint32_t get_table_offset(GlyphTableMap *glyph_map, const char *tag) {
  for (int i = 0; i < glyph_map->count; ++i) {
    if (strncmp(glyph_map->tables[i].tag, tag, 4) == 0) {
      return glyph_map->tables[i].offset;
    }
  }
  nob_log(NOB_ERROR, "Table tag '%s' not found", tag);
  return 0;
}
