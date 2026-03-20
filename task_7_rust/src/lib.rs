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

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_multiply_by_two_positive() {
        assert_eq!(multiply_by_two(5).unwrap(), 10);
        assert_eq!(multiply_by_two(100).unwrap(), 200);
    }

    #[test]
    fn test_multiply_by_two_zero() {
        assert_eq!(multiply_by_two(0).unwrap(), 0);
    }

    #[test]
    fn test_multiply_by_two_negative() {
        assert_eq!(multiply_by_two(-5).unwrap(), -10);
        assert_eq!(multiply_by_two(-100).unwrap(), -200);
    }

    #[test]
    fn test_make_greeting() {
        assert_eq!(
            make_greeting("Alice").unwrap(),
            "Hi, Alice! Welcome from Rust."
        );
        assert_eq!(
            make_greeting("Bob").unwrap(),
            "Hi, Bob! Welcome from Rust."
        );
    }
}