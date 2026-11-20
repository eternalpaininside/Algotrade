#pragma once
#include <cstddef>
#include <vector>
#include <algorithm>
#include <cmath>
#include <optional>

using std::abs;
using std::vector;

struct ForecastResult {
  size_t target_index;
  double predicted_price;
  double actual_price;
  double error;
};

vector<double> Calculate_SMA(const vector<double>& prices, const size_t window);

std::optional<double> Calculate_MAE(const vector<ForecastResult>& results,
                                    size_t first_target_index);

vector<ForecastResult> Calculate_Momentum(const vector<double>& prices, const size_t window);

vector<double> Calculate_EMA(const vector<double>& prices, const double alpha);

vector<double> Calculate_linear_trend(const vector<double>& prices, const size_t window);