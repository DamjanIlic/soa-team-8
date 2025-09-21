import { Routes } from '@angular/router';
import { UserListComponent } from './features/admin/user-list/user-list.component';
import { BlogListComponent } from './features/blog/blog-list/blog-list.component';
import { BlogDetailsComponent } from './features/blog/blog-details/blog-details.component';


export const routes: Routes = [

    { path: 'admin/users', component: UserListComponent },
    { 
      path: 'blogs', 
      loadComponent: () => import('./features/blog/blog-list/blog-list.component').then(m => m.BlogListComponent)
    },
    { path: 'blogs/:id', component: BlogDetailsComponent }
];
