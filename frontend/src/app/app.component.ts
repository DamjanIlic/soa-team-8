import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterModule, RouterOutlet } from "@angular/router";
import { NavbarComponent } from './shared/components/navbar/navbar.component';
import { AuthService } from './core/services/auth.service';
import { ShoppingCartComponent } from './features/cart/shopping-cart.component';

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    standalone: true,
    imports: [RouterOutlet ,CommonModule, RouterModule, NavbarComponent, ShoppingCartComponent]
})
export class AppComponent {
    constructor(private authService: AuthService) {}
}
