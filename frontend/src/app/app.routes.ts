import { Routes } from '@angular/router';
import { RegisterComponent } from './features/auth/register/register.component';
import { LoginComponent } from './features/auth/login/login.component';
import { WelcomeComponent } from './features/auth/welcome/welcome.component';
import { DashboardComponent } from './features/stakeholders/dashboard/dashboard.component';
import { ProfileComponent } from './features/stakeholders/profile/profile.component';

export const routes: Routes = [
  { path: '', component: WelcomeComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'login', component: LoginComponent },
  { 
    path: 'dashboard', component: DashboardComponent,
    children: [
      { path: 'profile', component: ProfileComponent }
    ]
  },
];
