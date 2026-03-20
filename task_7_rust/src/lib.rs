use pyo3::prelude::*;

#[pyfunction]
fn multiply_by_two(x: i32) -> PyResult<i32> {
    Ok(x * 2)
}

#[pyfunction]
fn make_greeting(name: &str) -> PyResult<String> {
    Ok(format!("Hi, {}! Welcome from Rust.", name))
}

#[pymodule]
fn my_rust_module(_py: Python<'_>, m: &Bound<'_, PyModule>) -> PyResult<()> {
    m.add_function(wrap_pyfunction!(multiply_by_two, m)?)?;
    m.add_function(wrap_pyfunction!(make_greeting, m)?)?;
    Ok(())
}