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
package com.iluwatar.masterworker;

import com.iluwatar.masterworker.system.ArrayTransposeMasterWorker;
import com.iluwatar.masterworker.system.systemmaster.ArrayTransposeMaster;
import com.iluwatar.masterworker.system.systemworkers.ArrayTransposeWorker;
import lombok.extern.slf4j.Slf4j;

/**
 * The <b><em>Master-Worker</em></b> pattern is used when the problem at hand can be solved by
 * dividing into multiple parts which need to go through the same computation and may need to be
 * aggregated to get final result. Parallel processing is performed using a system consisting of a
 * master and some number of workers, where a master divides the work among the workers, gets the
 * result back from them and assimilates all the results to give final result. The only
 * communication is between the master and the worker - none of the workers communicate among one
 * another and the user only communicates with the master to get required job done.
 *
 * <p>In our example, {@link ArrayTransposeMasterWorker} is the entry point used by clients. It
 * holds an {@link ArrayTransposeMaster}, which divides the input matrix into row-chunks, starts
 * one {@link ArrayTransposeWorker} thread per chunk, waits for them to finish, and reassembles
 * their partial results into the final transposed matrix. {@link ArrayInput} and {@link
 * ArrayResult} are simple wrappers around the matrix data passed between master and workers.
 */
@Slf4j
public class App {

  /**
   * Program entry point.
   *
   * @param args command line args
   */
  public static void main(String[] args) {
    var mw = new ArrayTransposeMasterWorker();
    var rows = 10;
    var columns = 20;
    var inputMatrix = ArrayUtilityMethods.createRandomIntMatrix(rows, columns);
    var input = new ArrayInput(inputMatrix);
    var result = (ArrayResult) mw.getResult(input);
    if (result != null) {
      ArrayUtilityMethods.printMatrix(inputMatrix);
      ArrayUtilityMethods.printMatrix(result.data);
    } else {
      LOGGER.info("Please enter non-zero input");
    }
  }
}
