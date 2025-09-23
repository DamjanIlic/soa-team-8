import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CartService, Cart, CartItem } from '../../core/services/cart.service';
import { Subscription } from 'rxjs';
import { Router } from '@angular/router';

@Component({
  selector: 'app-shopping-cart',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './shopping-cart.component.html',
  styleUrls: ['./shopping-cart.component.css']
})
export class ShoppingCartComponent implements OnInit, OnDestroy {
  isVisible = false;
  cart: Cart | null = null;
  loading = false;
  error = '';
  
  private subscriptions = new Subscription();

  constructor(
    private cartService: CartService,
    private router: Router
  ) {}

  ngOnInit(): void {
    // Subscribe to cart visibility
    this.subscriptions.add(
      this.cartService.cartVisible$.subscribe(visible => {
        this.isVisible = visible;
      })
    );

    // Subscribe to cart data
    this.subscriptions.add(
      this.cartService.cart$.subscribe(cart => {
        this.cart = cart;
      })
    );
  }

  ngOnDestroy(): void {
    this.subscriptions.unsubscribe();
  }

  closeCart(): void {
    this.cartService.hideCart();
  }

  removeItem(itemId: string): void {
    this.loading = true;
    this.cartService.removeItem(itemId).subscribe({
      next: () => {
        this.cartService.refreshCart();
        this.loading = false;
      },
      error: (err) => {
        console.error('Error removing item:', err);
        this.error = 'Failed to remove item';
        this.loading = false;
      }
    });
  }

  checkout(): void {
    if (!this.cart || this.cart.items.length === 0) {
      return;
    }

    this.loading = true;
    this.cartService.checkout().subscribe({
      next: (tokens) => {
        alert('Purchase successful! You have received tour tokens.');
        this.cartService.refreshCart();
        this.cartService.hideCart();
        this.loading = false;

        this.router.navigate(['/tours/my-tours']);
      },
      error: (err) => {
        console.error('Checkout error:', err);
        this.error = 'Checkout failed. Please try again.';
        this.loading = false;
      }
    });
  }

  getTotalItems(): number {
    return this.cart?.items.length || 0;
  }

  formatPrice(price: number): string {
    return `$${price.toFixed(2)}`;
  }
}