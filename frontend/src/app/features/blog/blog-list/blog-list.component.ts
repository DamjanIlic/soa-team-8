import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BlogService } from '../../../core/services/blog.service'; 
import { Blog } from '../../../core/models/blog.model'; 
import { catchError, of } from 'rxjs';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-blog-list',
  standalone: true,
  imports: [CommonModule, RouterModule], // HttpClientModule više nije potrebno
  templateUrl: './blog-list.component.html',
  styleUrls: ['./blog-list.component.css']
})
export class BlogListComponent implements OnInit {
  blogs: Blog[] = [];
  blogsForUser: Blog[] = [];
  loading = false;
  error = '';
  selectedTab: 'foryou' | 'explore' = 'foryou';

  constructor(private blogService: BlogService) {}

  ngOnInit(): void {
    this.loadForYouBlogs();
    //this.fetchBlogs();
  }

  fetchBlogs() {
    this.loading = true;
    this.blogService.getAll()
      .pipe(
        catchError(err => {
          this.error = 'Error fetching blogs.';
          console.error(err);
          return of([]);
        })
      )
      .subscribe(data => {
        this.blogs = data;
        this.loading = false;
      });
  }

  loadForYouBlogs() {
    this.selectedTab = 'foryou';
    this.loading = true;
    this.blogService.getForCurrentUser() // nova funkcija u blogService
      .pipe(
        catchError(err => {
          this.error = 'Error fetching For You blogs.';
          console.error(err);
          return of([]);
        })
      )
      .subscribe(data => {
        this.blogsForUser = data;
        this.blogs = this.blogsForUser;
        this.loading = false;
      });
  }

  loadExploreBlogs() {
    this.selectedTab = 'explore';
    this.fetchBlogs(); // koristi postojeću getAll funkciju
  }
}
