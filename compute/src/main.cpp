#include <algorithm>
#include <cmath>
#include <iomanip>
#include <iostream>
#include <optional>
#include <vector>

using std::abs;
using std::cin;
using std::cout;
using std::endl;
using std::vector;

struct ForecastResult {
  size_t target_index;
  double predicted_price;
  double actual_price;
  double error;
};

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

int main() {
  setlocale(LC_ALL, "RU");
  std::ios_base::sync_with_stdio(false);
  cin.tie(nullptr);

  vector<double> prices{
      100, 102, 101, 104, 105, 100, 101, 101, 103, 101.5
  };

  vector<ForecastResult> results;

  if (prices.size() >= 2) {
    for (size_t i = 1; i < prices.size(); ++i) {
      results.push_back({
          i,
          prices[i - 1],
          prices[i],
          prices[i - 1] - prices[i]
        });
    }
  } else {
    std::cerr << "Данных недостаточно\n";
    return 1;
  }

  cout << "------------\n";

  for (const auto& i : prices) {
    cout << i << " ";
  }

  cout << "\n------------\n";

  const vector<double> sma_predictions = Calculate_SMA(prices);
  vector<ForecastResult> sma_results;

  for (size_t i = 3; i < prices.size(); ++i) {
    const double sma_price = sma_predictions[i - 3];

    sma_results.push_back({
        i,
        sma_price,
        prices[i],
        sma_price - prices[i]
      });
  }

  cout << std::left << std::setw(12) << "Target"
    << std::right << std::setw(12) << "Actual"
    << std::setw(16) << "Last price"
    << std::setw(12) << "SMA"
    << std::setw(16) << "Last Error"
    << std::setw(12) << "SMA error" << '\n';

  cout << std::fixed << std::setprecision(2);

  for (size_t i = 0; i < prices.size(); ++i) {
    cout << std::left << std::setw(12) << i
      << std::right << std::setw(12) << prices[i];

    if (i >= 1) {
      cout << std::setw(16) << results[i - 1].predicted_price;
    } else {
      cout << std::setw(16) << "";
    }

    if (i >= 3) {
      cout << std::setw(12) << sma_results[i - 3].predicted_price;
    } else {
      cout << std::setw(12) << "";
    }

    if (i >= 1) {
      cout << std::setw(16) << results[i - 1].error;
    } else {
      cout << std::setw(16) << "";
    }

    if (i >= 3) {
      cout << std::setw(12) << sma_results[i - 3].error;
    } else {
      cout << std::setw(12) << "";
    }

    cout << '\n';
  }

  const auto last_price_mae = Calculate_MAE(results, 3);
  const auto sma_mae = Calculate_MAE(sma_results, 3);

  if (last_price_mae.has_value()) {
    cout << "MAE Last price = " << *last_price_mae << '\n';
  } else {
    cout << "MAE Last price посчитать невозможно\n";
  }

  if (sma_mae.has_value()) {
    cout << "MAE SMA = " << *sma_mae << '\n';
  } else {
    cout << "MAE SMA посчитать невозможно\n";
  }

  cout << "\n------------\n";
}