#pragma once
#include <errno.h>
#include <math.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <time.h>

#include "lib/nob.h"
#include "lib/stb_ds.h"

#include "lib/raylib/raylib-5.5_linux_amd64/include/raylib.h"
#include "lib/raylib/raylib-5.5_linux_amd64/include/rlgl.h"

typedef struct {
  char tag[5];
  __uint32_t offset;
} TagOffsetEntry;

typedef struct {
  TagOffsetEntry *tables;
  int count;
} TagOffsetMap;

typedef struct {
  __uint32_t unicode;
  __uint32_t index;
} GlyphUnicodeIndexEntry;

typedef struct {
  GlyphUnicodeIndexEntry *indices;
  int count;
} GlyphUnicodeIndexMap;

typedef struct {
  long read_loc;
  __uint16_t offset;
} IdRangeOffset;

// GLYPH
__uint32_t *GetAllGlyphLocations(FILE *f, TagOffsetMap *tag_offset_map);
GlyphUnicodeIndexMap *GetUnicodeGlyphIndexMap(FILE *f, TagOffsetMap *tag_map);

__uint32_t find_best_cmap_subtable(FILE *f, __uint16_t num_subtables);
GlyphUnicodeIndexMap *parse_format_12_cmap(FILE *f);
GlyphUnicodeIndexMap *parse_format_4_cmap(FILE *f);

// UTILS
long get_loc(FILE *file);
void goto_loc(FILE *file, long position);
void skip_bytes(FILE *file, long num);

__uint8_t read_uint8(FILE *file);
__int16_t read_int16(FILE *file);
__uint16_t read_uint16(FILE *file);
__uint32_t read_uint32(FILE *file);
void read_tag(FILE *file, char *tag);
IdRangeOffset read_id_range_offset(FILE *f);

bool flag_bit_is_set(__uint32_t flag, int bit_index);
__uint32_t get_tag_offset(TagOffsetMap *tag_offset_map, const char *tag);

// MACROS
#define READ_INTO_ARRAY(dst_array_ptr, src_file, type, read_fn, count)         \
  do {                                                                         \
    (dst_array_ptr) = calloc((count), sizeof(type));                           \
    if ((!dst_array_ptr)) {                                                    \
      nob_log(NOB_ERROR, "calloc failed for '" #dst_array_ptr "'");            \
      exit(1);                                                                 \
    }                                                                          \
    for (int _i = 0; _i < (count); ++_i) {                                     \
      (dst_array_ptr)[_i] = read_fn(src_file);                                 \
    }                                                                          \
  } while (0)\
