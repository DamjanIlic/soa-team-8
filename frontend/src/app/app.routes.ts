import { Routes } from '@angular/router';
import { UserListComponent } from './features/admin/user-list/user-list.component';
import { BlogListComponent } from './features/blog/blog-list/blog-list.component';



export const routes: Routes = [

    { path: 'admin/users', component: UserListComponent },
    { path: 'blogs', component: BlogListComponent}
];
