/**
 * Copyright 2019 The Gamma Authors.
 *
 * This source code is licensed under the Apache License, Version 2.0 license
 * found in the LICENSE file in the root directory of this source tree.
 */

#pragma once

#include <tbb/concurrent_queue.h>

#include <map>
#include <string>
#include <vector>

#include "concurrentqueue/blockingconcurrentqueue.h"
#include "range_query_result.h"
#include "table.h"

namespace vearch {

enum class FilterOperator : uint8_t { And = 0, Or, Not };

typedef struct {
  int field;
  std::string lower_value;
  std::string upper_value;
  bool include_lower;
  bool include_upper;
  FilterOperator is_union;
} FilterInfo;

class FieldOperate {
 public:
  typedef enum { ADD, DELETE } operate_type;
  explicit FieldOperate(operate_type type, int doc_id, int field_id)
      : type(type), doc_id(doc_id), field_id(field_id) {}

  operate_type type;
  int doc_id;
  int field_id;
  std::string value;
};

typedef tbb::concurrent_bounded_queue<FieldOperate *> FieldOperateQueue;

class FieldRangeIndex;
class MultiFieldsRangeIndex {
 public:
  MultiFieldsRangeIndex(std::string &path, Table *table);
  ~MultiFieldsRangeIndex();

  int Add(int docid, int field);

  int Delete(int docid, int field);

  int AddField(int field, enum DataType field_type, std::string &name);

  int Search(const std::vector<FilterInfo> &origin_filters,
             MultiRangeQueryResults *out);

  // for debug
  long MemorySize(long &dense, long &sparse);

  int PendingTasks() { return field_operate_q_->size(); }

 private:
  int Intersect(std::vector<RangeQueryResult> &results, int shortest_idx,
                RangeQueryResult *out);
  void FieldOperateWorker();

  int AddDoc(int docid, int field);

  int DeleteDoc(int docid, int field, std::string &key);
  std::vector<FieldRangeIndex *> fields_;
  Table *table_;
  std::string path_;
  bool b_running_;
  bool b_operate_running_;
  FieldOperateQueue *field_operate_q_;
};

}  // namespace vearch
