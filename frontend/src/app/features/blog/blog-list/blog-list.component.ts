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
  loading = false;
  error = '';

  constructor(private blogService: BlogService) {}

  ngOnInit(): void {
    this.fetchBlogs();
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
}
