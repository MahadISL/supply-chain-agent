package com.supplychain.core_service;

import com.supplychain.core_service.model.Product;
import com.supplychain.core_service.model.Supplier;
import com.supplychain.core_service.repository.ProductRepository;
import com.supplychain.core_service.repository.SupplierRepository;
import org.springframework.boot.CommandLineRunner;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Component
public class DataSeeder implements CommandLineRunner {

    private final SupplierRepository supplierRepository;
    private final ProductRepository productRepository;

    public DataSeeder(SupplierRepository supplierRepository, ProductRepository productRepository) {
        this.supplierRepository = supplierRepository;
        this.productRepository = productRepository;
    }

    @Override
    @Transactional
    public void run(String... args) throws Exception {

        // 1. Check if data already exists to prevent duplicates
        if (supplierRepository.count() == 0) {
            System.out.println("--- STARTING DATA SEEDING ---");

            // 2. Create Suppliers
            Supplier supplier1 = new Supplier();
            supplier1.setName("Apex Furniture");
            supplier1.setEmail("apex_orders@example.com");
            supplier1.setContactInfo("123 Industrial Blvd, NY");

            Supplier supplier2 = new Supplier();
            supplier2.setName("TechGadget Inc");
            supplier2.setEmail("sales@techgadget.com");
            supplier2.setContactInfo("456 Silicon Ave, CA");

            // Save Suppliers first (so they generate IDs)
            supplierRepository.saveAll(List.of(supplier1, supplier2));

            // 3. Create Products linked to Suppliers

            // Product A: Healthy Stock (No Action needed)
            Product chair = new Product();
            chair.setName("Ergonomic Office Chair");
            chair.setSku("FURN-001");
            chair.setPrice(150.00);
            chair.setDescription("Mesh back support, adjustable height.");
            chair.setStockQuantity(50);
            chair.setMinStockLevel(10);
            chair.setSupplier(supplier1); // Linked to Apex

            // Product B: LOW STOCK (This will trigger the Agent later!)
            Product desk = new Product();
            desk.setName("Standing Desk Pro");
            desk.setSku("FURN-002");
            desk.setPrice(450.00);
            desk.setDescription("Dual motor electric standing desk.");
            desk.setStockQuantity(5);   // <--- DANGER! Below min level
            desk.setMinStockLevel(10);  // Trigger point
            desk.setSupplier(supplier1); // Linked to Apex

            // Product C: Healthy Stock
            Product monitor = new Product();
            monitor.setName("4K Monitor 27-inch");
            monitor.setSku("TECH-001");
            monitor.setPrice(299.99);
            monitor.setDescription("IPS Panel, 144Hz refresh rate.");
            monitor.setStockQuantity(100);
            monitor.setMinStockLevel(15);
            monitor.setSupplier(supplier2); // Linked to TechGadget

            // Save Products
            productRepository.saveAll(List.of(chair, desk, monitor));

            System.out.println("--- DATA SEEDING COMPLETED: 2 Suppliers, 3 Products ---");
        } else {
            System.out.println("--- DATA ALREADY EXISTS. SKIPPING SEEDING ---");
        }
    }
}