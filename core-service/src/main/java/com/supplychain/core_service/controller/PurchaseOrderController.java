package com.supplychain.core_service.controller;

import com.supplychain.core_service.model.OrderStatus;
import com.supplychain.core_service.model.Product;
import com.supplychain.core_service.model.PurchaseOrder;
import com.supplychain.core_service.repository.ProductRepository;
import com.supplychain.core_service.repository.PurchaseOrderRepository;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/orders")
@CrossOrigin(origins = "*")
public class PurchaseOrderController {

    private final PurchaseOrderRepository orderRepository;
    private final ProductRepository productRepository;

    public PurchaseOrderController(PurchaseOrderRepository orderRepository, ProductRepository productRepository) {
        this.orderRepository = orderRepository;
        this.productRepository = productRepository;
    }

    @GetMapping
    public List<PurchaseOrder> getAllOrders() {
        return orderRepository.findAll();
    }

    @PostMapping
    public ResponseEntity<PurchaseOrder> createOrder(@RequestBody OrderRequest request) {
        return productRepository.findById(request.productId())
                .map(product -> {
                    PurchaseOrder order = new PurchaseOrder();
                    order.setProduct(product);
                    order.setQuantity(request.quantity());
                    order.setStatus(OrderStatus.PENDING_APPROVAL); // Default status
                    order.setTotalCost(product.getPrice() * request.quantity());

                    return ResponseEntity.ok(orderRepository.save(order));
                })
                .orElse(ResponseEntity.notFound().build());
    }

    public record OrderRequest(Long productId, Integer quantity) {}
}
