use pyo3::prelude::*;

#[pyclass]
#[derive(Debug, Clone)]
pub struct Point {
    #[pyo3(get, set)]
    pub x: f64,
    #[pyo3(get, set)]
    pub y: f64,
    #[pyo3(get, set)]
    pub on_curve: bool,
}

#[pymethods]
impl Point {
    #[new]
    pub fn new(x: f64, y: f64, on_curve: bool) -> Self {
        Self { x, y, on_curve }
    }
}
