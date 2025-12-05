package com.supplychain.core_service.model;

import jakarta.persistence.*;
import lombok.Data;

@Entity
@Table(name = "products")
@Data
public class Product {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false)
    private String name;

    @Column(unique = true, nullable = false)
    private String sku; // Stock Keeping Unit (e.g., "FURN-001")

    private Double price;

    @Column(columnDefinition = "TEXT")
    private String description;

    // INVENTORY MANAGEMENT FIELDS
    @Column(nullable = false)
    private Integer stockQuantity; // Current amount in warehouse

    @Column(nullable = false)
    private Integer minStockLevel; // The "Threshold". If stock < this, Agent wakes up.

    // RELATIONSHIPS
    @ManyToOne(fetch = FetchType.EAGER) // When we load Product, load Supplier too
    @JoinColumn(name = "supplier_id", nullable = false)
    private Supplier supplier;
}
