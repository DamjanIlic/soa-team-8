import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Router } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './navbar.component.html'
})
export class NavbarComponent implements OnInit, OnDestroy {
  isLoggedIn = false;
  userRole: string = '';
  private authSubscription: Subscription = new Subscription();

  constructor(
    private router: Router,
    private authService: AuthService
  ) {}

  ngOnInit(): void {
    // Slusaj promene auth status-a
    this.authSubscription = this.authService.authStatus$.subscribe(
      (isAuthenticated: boolean) => {
        this.isLoggedIn = isAuthenticated;
        if (isAuthenticated) {
          this.userRole = this.authService.getUserRole();
        } else {
          this.userRole = '';
        }
      }
    );
  }

  ngOnDestroy(): void {
    this.authSubscription.unsubscribe();
  }

  logout(): void {
    this.authService.logout();
    this.router.navigate(['/auth/login']);
  }
  
}