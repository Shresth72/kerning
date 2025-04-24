use pyo3::prelude::*;

use lazy_static::lazy_static;
use std::sync::Mutex;

use bezier::{draw_bezier, Point};

mod bezier;

lazy_static! {
    static ref GLYPHS: Mutex<Vec<Glyph>> = Mutex::new(vec![]);
}

#[derive()]
struct Glyph {
    is_space: bool,
    contours: Vec<Vec<Point>>,
}

pub fn store_glyphs(glyphs: Vec<(bool, Vec<Vec<(f64, f64)>>)>) {
    let parsed = glyphs
        .into_iter()
        .map(|(is_space, contours)| Glyph {
            is_space,
            contours: contours
                .into_iter()
                .map(|c| c.into_iter().map(|(x, y)| [x, y]).collect())
                .collect(),
        })
        .collect();

    let mut lock = GLYPHS.lock().unwrap();
    *lock = parsed;
}

pub fn draw_glyphs() {
    let glyphs = GLYPHS.lock().unwrap();
    let resolution = 100;
    let spacing = 700.0;
    let word_spacing = 500.0;
    let mut x_cursor = 0.0;

    for glyph in glyphs.iter() {
        if glyph.is_space {
            x_cursor += word_spacing;
            continue;
        }

        for contour in &glyph.contours {
            let shifted: Vec<Point> = contour //
                .iter()
                .map(|p| [p[0] + x_cursor, p[1]])
                .collect();
            let len = shifted.len();

            let mut i = 0;
            while i + 2 < len {
                let p0 = shifted[i];
                let p1 = shifted[i + 1];
                let p2 = shifted[i + 2];
                draw_bezier(&p0, &p1, &p2, resolution);
                i += 2;
            }
        }

        x_cursor += spacing;
    }
}

#[pyfunction]
fn accept_glyphs(glyphs: Vec<(bool, Vec<Vec<(f64, f64)>>)>) -> PyResult<()> {
    store_glyphs(glyphs);
    draw_glyphs();
    Ok(())
}

#[pymodule]
fn renderer(m: &Bound<'_, PyModule>) -> PyResult<()> {
    m.add_function(wrap_pyfunction!(accept_glyphs, m)?)?;
    Ok(())
}
