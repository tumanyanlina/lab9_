import pytest

try:
    import my_rust_module
except ImportError:
    my_rust_module = None


def test_module_import():
    """Проверка, что модуль my_rust_module успешно импортируется."""
    assert my_rust_module is not None
    assert hasattr(my_rust_module, "multiply_by_two")
    assert hasattr(my_rust_module, "make_greeting")


class TestMultiplyByTwo:
    """Тесты для функции multiply_by_two."""

    def test_positive_number(self):
        """Умножение положительного числа на 2."""
        assert my_rust_module.multiply_by_two(5) == 10
        assert my_rust_module.multiply_by_two(100) == 200

    def test_zero(self):
        """Умножение нуля на 2."""
        assert my_rust_module.multiply_by_two(0) == 0

    def test_negative_number(self):
        """Умножение отрицательного числа на 2."""
        assert my_rust_module.multiply_by_two(-5) == -10
        assert my_rust_module.multiply_by_two(-100) == -200


class TestMakeGreeting:
    """Тесты для функции make_greeting."""

    def test_simple_name(self):
        """Приветствие с простым именем."""
        assert my_rust_module.make_greeting("Alice") == "Hi, Alice! Welcome from Rust."

    def test_another_name(self):
        """Приветствие с другим именем."""
        assert my_rust_module.make_greeting("Bob") == "Hi, Bob! Welcome from Rust."

    def test_name_with_spaces(self):
        """Приветствие с именем, содержащим пробелы."""
        result = my_rust_module.make_greeting("John Doe")
        assert result == "Hi, John Doe! Welcome from Rust."

    def test_empty_string(self):
        """Приветствие с пустой строкой."""
        assert my_rust_module.make_greeting("") == "Hi, ! Welcome from Rust."


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
