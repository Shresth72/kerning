use pyo3::prelude::*;

mod bezier;
use bezier::Point;

#[derive()]
struct Glyph {
    is_space: bool,
    contours: Vec<Vec<Point>>,
}

#[pyfunction]
fn accept_glyphs(glyphs: Vec<(bool, Vec<Vec<(f64, f64)>>)>) -> PyResult<()> {
    let parsed_glyphs: Vec<Glyph> = glyphs
        .into_iter()
        .map(|(is_space, contours)| Glyph {
            is_space,
            contours: contours
                .into_iter()
                .map(|c| c.into_iter().map(|(x, y)| [x, y]).collect())
                .collect(),
        })
        .collect();

    println!("Received glyphs");
    for glyph in parsed_glyphs.iter() {
        if glyph.is_space {
            println!("(space)");
        } else {
            println!("Contours: {:?}", glyph.contours);
        }
    }
    Ok(())
}

#[pymodule]
fn renderer(m: &Bound<'_, PyModule>) -> PyResult<()> {
    m.add_function(wrap_pyfunction!(accept_glyphs, m)?)?;
    Ok(())
}
