package org.ivangalochkin.expensewriter;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@SpringBootApplication
public class ExpenseWriterApplication {
    private static final Logger log = LoggerFactory.getLogger(ExpenseWriterApplication.class);

    public static void main(String[] args) {
        log.info("Starting ExpenseWriterApplication");
        SpringApplication.run(ExpenseWriterApplication.class, args);
        log.info("ExpenseWriterApplication started");
    }
}
