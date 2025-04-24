use pyo3::prelude::*;

mod bezier;
use bezier::Point;

#[pyfunction]
fn accept_countours(contours: Vec<Vec<(f64, f64, bool)>>) -> PyResult<()> {
    let parsed_contours: Vec<Vec<Point>> = contours
        .into_iter()
        .map(|contour| {
            contour
                .into_iter()
                .map(|(x, y, on_curve)| Point { x, y, on_curve })
                .collect()
        })
        .collect();
    println!("Received contours: {:?}", parsed_contours);
    Ok(())
}

#[pymodule]
fn renderer(m: &Bound<'_, PyModule>) -> PyResult<()> {
    m.add_function(wrap_pyfunction!(accept_countours, m)?)?;
    Ok(())
}
