use pyo3::prelude::*;

/// Возвращает квадрат числа n.
#[pyfunction]
fn square(n: i32) -> i32 {
    n * n
}

/// Возвращает сумму квадратов элементов списка.
#[pyfunction]
fn sum_of_squares(numbers: Vec<i32>) -> i32 {
    numbers.iter().map(|&x| x * x).sum()
}

/// Модуль Python rust_math.
#[pymodule]
fn rust_math(_py: Python, m: &PyModule) -> PyResult<()> {
    m.add_function(wrap_pyfunction!(square, m)?)?;
    m.add_function(wrap_pyfunction!(sum_of_squares, m)?)?;
    Ok(())
}
