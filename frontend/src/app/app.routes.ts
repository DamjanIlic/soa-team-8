import { Routes } from '@angular/router';
import { UserListComponent } from './features/admin/user-list/user-list.component';
import { BlogListComponent } from './features/blog/blog-list/blog-list.component';
import { BlogDetailsComponent } from './features/blog/blog-details/blog-details.component';
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
    { path: 'admin/users', component: UserListComponent },
    { 
        path: 'blogs', 
        loadComponent: () => import('./features/blog/blog-list/blog-list.component').then(m => m.BlogListComponent)
    },
    { path: 'blogs/:id', component: BlogDetailsComponent }
];
