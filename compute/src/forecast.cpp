#include "forecast.h"

vector<double> Calculate_SMA(const vector<double>& prices) {
  if (prices.size() < 4) {
    return vector<double>{};
  }

  vector<double> sma;

  for (size_t i = 3; i < prices.size(); ++i) {
    sma.push_back(
      (prices[i - 2] + prices[i - 1] + prices[i - 3]) / 3
    );
  }

  return sma;
}

std::optional<double> Calculate_MAE(
  const vector<ForecastResult>& results,
  size_t first_target_index
) {
  double total_error = 0;
  size_t count = 0;

  for (size_t i = 0; i < results.size(); ++i) {
    if (results[i].target_index < first_target_index) {
      continue;
    }

    total_error += std::abs(results[i].error);
    ++count;
  }

  if (count == 0) {
    return std::nullopt;
  }

  double mae = total_error / static_cast<double>(count);
  return mae;
}