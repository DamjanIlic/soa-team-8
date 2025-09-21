import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { MarkdownModule } from 'ngx-markdown';
import { BlogService } from '../../../core/services/blog.service';
import { Blog } from '../../../core/models/blog.model';

@Component({
  selector: 'app-blog-details',
  standalone: true,
  imports: [CommonModule, MarkdownModule],
  templateUrl: './blog-details.component.html',
  styleUrl: './blog-details.component.css'
})
export class BlogDetailsComponent implements OnInit {
  blog?: Blog;
  blogId!: string;
  //tmp jwt
  token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTg3MzEwOTIsInJvbGUiOiJndWlkZSIsInVzZXJfaWQiOiIzMzUyODRhZC01YTBiLTRmOTAtYjUzNC1iOWFlZjdiYjVkNDMifQ.8nW9EVjIQR6et5BgLKkHPQ_s9AM01LYu_nG8-wgRWso"
  constructor(
    private route: ActivatedRoute,
    private blogService: BlogService
  ) {}

  ngOnInit(): void {
    this.blogId = this.route.snapshot.paramMap.get('id')!;
    this.blogService.getById(this.blogId).subscribe({
      next: (data) => this.blog = data,
      error: (err) => console.error(err)
    });
  }

  likeBlog(): void {
    this.blogService.like(this.blogId, this.token).subscribe({
      next: (likes) => {
        if (this.blog) this.blog.likes = likes;
      },
      error: (err) => console.error(err)
    });
  }

  unlikeBlog(): void {
    this.blogService.unlike(this.blogId, this.token).subscribe({
      next: (likes) => {
        if (this.blog) this.blog.likes = likes;
      },
      error: (err) => console.error(err)
    });
  }
}
