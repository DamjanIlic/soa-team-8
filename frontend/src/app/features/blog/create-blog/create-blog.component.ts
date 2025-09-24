import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { BlogRequest } from '../../../core/models/blog.model';
import { BlogService } from '../../../core/services/blog.service';
import { MarkdownModule } from 'ngx-markdown';

@Component({
  selector: 'app-create-blog',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule, MarkdownModule],
  templateUrl: './create-blog.component.html'
})
export class CreateBlogComponent {
  blog: BlogRequest = {
    title: '',
    content: '',
    image_url: ''
  };

  message = '';

  constructor(private blogService: BlogService, private router: Router) {}

  createBlog() {
    if (!this.blog.title || !this.blog.content) {
      this.message = 'Title and Content are required!';
      return;
    }

    this.blogService.create(this.blog).subscribe({
      next: (res) => {
        this.router.navigate(['/blogs']);
      },
      error: (err) => {
        this.message = 'Failed to create blog. ' + (err.error?.message || '');
      }
    });
  }
}
