import { Routes } from '@angular/router';
import { UserListComponent } from './features/admin/user-list/user-list.component';
import { LoginComponent } from './features/auth/login/login.component';
import { RegisterComponent } from './features/auth/register/register.component';
import { WelcomeComponent } from './features/auth/welcome/welcome.component';
import { BlogDetailsComponent } from './features/blog/blog-details/blog-details.component';
import { DashboardComponent } from './features/stakeholders/dashboard/dashboard.component';
import { ProfileComponent } from './features/stakeholders/profile/profile.component';


import { BlogFeedComponent } from './features/blog/blog-feed/blog-feed.component';
import { CreateBlogComponent } from './features/blog/create-blog/create-blog.component';

import { MyToursComponent } from './features/tours/my-tours/my-tours.component';
import { AuthorToursComponent } from './features/tours/author-tours/author-tours.component';
import { PositionSimulatorComponent } from './features/tours/position-simulator/position-simulator.component';


import { ReviewsComponent } from './features/reviews/reviews.component';


import { ShoppingCartComponent } from './features/cart/shopping-cart.component';
import { BrowseToursComponent } from './features/tours/browse-tours/browse-tours.component';
import { TourExecutionComponent } from './features/tours/tour-execution/tour-execution.component';

export const routes: Routes = [
  // Auth routes
  { path: '', component: WelcomeComponent },
  { path: 'auth/register', component: RegisterComponent },
  { path: 'auth/login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'login', component: LoginComponent },

    // Profile (direktno, bez dashboard wrapper-a)
  { path: 'profile', component: ProfileComponent },
  
  // Admin routes
  // { path: 'admin/users', component: UsersComponent },
  
  // Blog routes
  { path: 'blog/create', component: CreateBlogComponent },
  { path: 'blog/feed', component: BlogFeedComponent },
  
  // Tours routes
  { path: 'tours/browse', component: BrowseToursComponent },
  {
  path: 'tours/create',
  loadComponent: () =>
    import('./features/tours/create-tour/create-tour.component')
      .then(m => m.CreateTourComponent)
},
{
  path: 'tours/add-tour-checkpoints',
  loadComponent: () =>
    import('./features/tours/add-tour-checkpoints/add-tour-checkpoints.component')
      .then(m => m.AddTourCheckpointsComponent)
},

  { path: 'tours/my-tours', component: MyToursComponent },
  { path: 'tours/author-tours', component: AuthorToursComponent },
  { path: 'tours/position-simulator', component: PositionSimulatorComponent },
  { path: 'tours/tour-execution/:id', component: TourExecutionComponent},
  // Reviews (jedna komponenta za sve uloge)
  { path: 'reviews', component: ReviewsComponent },
  
  // Cart
  { path: 'cart', component: ShoppingCartComponent },
  
  // Dashboard
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
