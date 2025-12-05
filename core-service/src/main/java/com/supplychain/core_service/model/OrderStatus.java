package com.supplychain.core_service.model;

public enum OrderStatus {
    PENDING_APPROVAL,   // 1. Agent drafted the order, waiting for human
    APPROVED,           // 2. Human Manager clicked "Approve"
    REJECTED,           // 3. Human Manager clicked "Reject"
    SENT_TO_SUPPLIER    // 4. Email successfully sent via AWS SES
}