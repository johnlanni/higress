#pragma once

#include <stdexcept>
#include <string>

#include "re2/re2.h"

namespace Wasm::Common::Regex {

class CompiledGoogleReMatcher {
 public:
  CompiledGoogleReMatcher(const std::string& regex,
                          bool do_program_size_check = true)
      : regex_(regex, re2::RE2::Quiet) {
    if (!regex_.ok()) {
      throw std::runtime_error(regex_.error());
    }
    if (do_program_size_check) {
      const auto regex_program_size =
          static_cast<uint32_t>(regex_.ProgramSize());
      if (regex_program_size > 100) {
        throw std::runtime_error("too complex regex: " + regex);
      }
    }
  }

  bool match(std::string_view value) const {
    return re2::RE2::FullMatch(re2::StringPiece(value.data(), value.size()),
                               regex_);
  }

  std::string replaceAll(std::string_view value,
                         std::string_view substitution) const {
    std::string result = std::string(value);
    re2::RE2::GlobalReplace(
        &result, regex_,
        re2::StringPiece(substitution.data(), substitution.size()));
    return result;
  }

 private:
  const re2::RE2 regex_;
};

}  // namespace Wasm::Common::Regex
