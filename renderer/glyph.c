#include "utils.h"

__uint32_t *GetAllGlyphLocations(FILE *f, TagOffsetMap *tag_offset_map) {
  goto_loc(f, get_tag_offset(tag_offset_map, "maxp") + 4);
  __uint16_t num_glyphs = read_uint16(f);

  goto_loc(f, get_tag_offset(tag_offset_map, "head"));
  skip_bytes(f, 50);

  bool is_two_byte_entry = (read_uint16(f) == 0);

  __uint32_t loc_table_start = get_tag_offset(tag_offset_map, "loca");
  __uint32_t glyph_table_start = get_tag_offset(tag_offset_map, "glyf");

  __uint32_t *all_glyph_locs = calloc(num_glyphs, sizeof(__uint32_t));
  if (!all_glyph_locs) {
    nob_log(NOB_ERROR, "calloc failed for 'all_glyph_locs'");
    exit(1);
  }

  for (int glyph_index = 0; glyph_index < num_glyphs; ++glyph_index) {
    __uint32_t glyph_data_pos =
        loc_table_start + glyph_index * (is_two_byte_entry ? 2 : 4);
    goto_loc(f, glyph_data_pos);

    __uint32_t glyph_data_offset =
        is_two_byte_entry ? read_uint16(f) * 2 : read_uint32(f);
    all_glyph_locs[glyph_index] = glyph_table_start + glyph_data_offset;
  }

  return all_glyph_locs;
}

GlyphUnicodeIndexMap *GetUnicodeGlyphIndexMap(FILE *f, TagOffsetMap *tag_map) {
  __uint32_t cmap_offset = get_tag_offset(tag_map, "cmap");
  goto_loc(f, cmap_offset);

  read_uint16(f);
  __uint16_t num_subtables = read_uint16(f);
  __uint32_t cmap_subtable_offset = find_best_cmap_subtable(f, num_subtables);

  if (cmap_subtable_offset == 0) {
    NOB_TODO("Font does not contain supported character map type");
  }
  goto_loc(f, cmap_offset + cmap_subtable_offset);

  __uint16_t format = read_uint16(f);
  if (format == 12) {
    return parse_format_12_cmap(f);
  } else if (format == 4) {
    return parse_format_12_cmap(f);
  } else {
    NOB_TODO("Font cmap format not supported yet");
  }
}

__uint32_t find_best_cmap_subtable(FILE *f, __uint16_t num_subtables) {
  __uint32_t best_offset = UINT32_MAX;

  for (__uint16_t i = 0; i < num_subtables; ++i) {
    __uint16_t platform_id = read_uint16(f);
    __uint16_t platform_specific_id = read_uint16(f);
    __uint32_t offset = read_uint32(f);

    if (platform_id == 0) {
      if (platform_specific_id == 4) {
        best_offset = offset;
        break;
      }
      if (platform_specific_id == 3 && best_offset == UINT32_MAX) {
        best_offset = offset;
      }
    }
  }

  return best_offset;
}

GlyphUnicodeIndexMap *parse_format_12_cmap(FILE *f) {
  GlyphUnicodeIndexMap *unicode_index_map =
      malloc(sizeof(GlyphUnicodeIndexMap));
  if (!unicode_index_map) {
    nob_log(NOB_ERROR, "malloc failed for 'unicode_index_map'");
    exit(1);
  }

  read_uint16(f);
  read_uint32(f);
  __uint32_t num_groups = read_uint32(f);

  unicode_index_map->indices =
      malloc(sizeof(GlyphUnicodeIndexEntry) * num_groups);
  if (!unicode_index_map->indices) {
    nob_log(NOB_ERROR, "malloc failed for 'unicode_index_map->indices'");
    exit(1);
  }

  for (__uint32_t i = 0; i < num_groups; ++i) {
    __uint32_t start_char_code = read_uint32(f);
    __uint32_t end_char_code = read_uint32(f);
    __uint32_t start_glyph_index = read_uint32(f);

    __uint32_t num_chars = end_char_code - start_char_code + 1;
    for (__uint32_t char_code = start_char_code; char_code <= end_char_code;
         ++char_code) {
      unicode_index_map->indices[i].unicode = char_code;
      unicode_index_map->indices[i].index =
          start_glyph_index + (char_code - start_char_code);
      i++;
    }
  }

  unicode_index_map->count = num_groups;
  return unicode_index_map;
}

GlyphUnicodeIndexMap *parse_format_4_cmap(FILE *f) {
  GlyphUnicodeIndexMap *unicode_index_map =
      malloc(sizeof(GlyphUnicodeIndexMap));
  if (!unicode_index_map) {
    nob_log(NOB_ERROR, "malloc failed for 'unicode_index_map'");
    exit(1);
  }

  skip_bytes(f, 4);
  __uint16_t seg_count_2x = read_uint16(f);
  __uint16_t seg_count = seg_count_2x / 2;
  skip_bytes(f, 6);

  __uint16_t *end_codes;
  READ_INTO_ARRAY(end_codes, f, __uint16_t, read_uint16, seg_count);

  skip_bytes(f, 2);

  __uint16_t *start_codes;
  READ_INTO_ARRAY(start_codes, f, __uint16_t, read_uint16, seg_count);

  __uint16_t *id_deltas;
  READ_INTO_ARRAY(id_deltas, f, __uint16_t, read_uint16, seg_count);

  IdRangeOffset *id_range_offsets;
  READ_INTO_ARRAY(id_range_offsets, f, IdRangeOffset, read_id_range_offset,
                  seg_count);

  NOB_TODO("TODO parse_4_format");
}
