import unittest
from client import calculate_squares

class TestCalculateSquares(unittest.TestCase):
    def test_normal_input(self):
        numbers = [1, 2, 3]
        result = calculate_squares(numbers)
        self.assertIsNotNone(result, "Result should not be None")
        self.assertEqual(result.get('sum'), 1*1 + 2*2 + 3*3)
        self.assertEqual(result.get('original'), numbers)
        self.assertFalse('error' in result and result['error'], f"Unexpected error: {result.get('error')}")

    def test_empty_array(self):
        numbers = []
        result = calculate_squares(numbers)
        self.assertIsNotNone(result, "Result should not be None")
        self.assertEqual(result.get('error'), "no numbers provided")
        self.assertTrue('sum' not in result or result.get('sum') is None or result.get('sum') == 0)
        self.assertTrue('original' not in result or result.get('original') is None or result.get('original') == [] or result.get('original') == '')

    def test_number_too_large(self):
        numbers = [10, 2000, 2]
        result = calculate_squares(numbers)
        self.assertIsNotNone(result, "Result should not be None")
        self.assertEqual(result.get('error'), "number too large")
        self.assertTrue('sum' not in result or result.get('sum') is None or result.get('sum') == 0)
        self.assertTrue('original' not in result or result.get('original') is None or result.get('original') == [] or result.get('original') == '')

if __name__ == "__main__":
    unittest.main()