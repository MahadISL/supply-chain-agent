package com.supplychain.core_service.repository;

import com.supplychain.core_service.model.Product;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import java.util.List;

@Repository
public interface ProductRepository extends JpaRepository<Product, Long> {

    List<Product> findByStockQuantityLessThanEqual(Integer minStockLevel);
}