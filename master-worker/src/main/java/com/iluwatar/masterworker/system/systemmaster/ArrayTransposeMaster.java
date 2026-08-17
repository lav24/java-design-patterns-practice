/*
 * This project is licensed under the MIT license. Module model-view-viewmodel is using ZK framework licensed under LGPL (see lgpl-3.0.txt).
 *
 * The MIT License
 * Copyright © 2014-2022 Ilkka Seppälä
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */
package com.iluwatar.masterworker.system.systemmaster;

import com.iluwatar.masterworker.ArrayInput;
import com.iluwatar.masterworker.ArrayResult;
import com.iluwatar.masterworker.system.systemworkers.ArrayTransposeWorker;
import java.util.ArrayList;
import java.util.Hashtable;
import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.IntStream;
import lombok.Getter;

/**
 * Splits an input matrix across {@link ArrayTransposeWorker} threads, waits for all of them to
 * finish, then reassembles their partial results into the final transposed matrix.
 */
public class ArrayTransposeMaster {

  private final int numOfWorkers;
  private final List<ArrayTransposeWorker> workers;
  private final Hashtable<Integer, ArrayResult> allResultData;
  private int expectedNumResults;
  @Getter private ArrayResult finalResult;

  public ArrayTransposeMaster(int numOfWorkers) {
    this.numOfWorkers = numOfWorkers;
    this.workers =
        IntStream.range(0, numOfWorkers)
            // ith worker gets id i+1
            .mapToObj(i -> new ArrayTransposeWorker(this, i + 1))
            .collect(Collectors.toCollection(() -> new ArrayList<>(numOfWorkers)));
    this.allResultData = new Hashtable<>(numOfWorkers);
  }

  public void doWork(ArrayInput input) {
    var dividedInput = input.divideData(numOfWorkers);
    if (dividedInput == null) {
      return;
    }
    this.expectedNumResults = dividedInput.size();
    for (var i = 0; i < this.expectedNumResults; i++) {
      // ith division given to ith worker in this.workers
      this.workers.get(i).setReceivedData(this, dividedInput.get(i));
      this.workers.get(i).start();
    }
    for (var i = 0; i < this.expectedNumResults; i++) {
      try {
        this.workers.get(i).join();
      } catch (InterruptedException e) {
        System.err.println("Error while executing thread");
      }
    }
  }

  public void receiveData(ArrayResult data, ArrayTransposeWorker w) {
    this.allResultData.put(w.getWorkerId(), data);
    if (this.allResultData.size() == this.expectedNumResults) {
      // all data received
      this.finalResult = aggregateData();
    }
  }

  private ArrayResult aggregateData() {
    // number of rows in final result is number of rows in any of the obtained results
    var rows = this.allResultData.elements().nextElement().data.length;
    var columns = 0; // columns = sum of number of columns in all results obtained from workers
    for (var result : this.allResultData.values()) {
      columns += result.data[0].length;
    }
    var resultData = new int[rows][columns];
    var columnsDone = 0; // columns aggregated so far
    for (var i = 0; i < this.expectedNumResults; i++) {
      // result obtained from ith worker
      var workerId = this.workers.get(i).getWorkerId();
      var work = this.allResultData.get(workerId).data;
      for (var m = 0; m < work.length; m++) {
        // m = row number, n = column number
        System.arraycopy(work[m], 0, resultData[m], columnsDone, work[0].length);
      }
      columnsDone += work[0].length;
    }
    return new ArrayResult(resultData);
  }
}
