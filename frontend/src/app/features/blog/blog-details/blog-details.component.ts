import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { MarkdownModule } from 'ngx-markdown';
import { BlogService } from '../../../core/services/blog.service';
import { Blog } from '../../../core/models/blog.model';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule } from '@angular/forms'; // za ngModel
import { Comment } from '../../../core/models/comment.model';
@Component({
  selector: 'app-blog-details',
  standalone: true,
  imports: [
    CommonModule,
    HttpClientModule,
    MarkdownModule,
    FormsModule
  ],
  templateUrl: './blog-details.component.html',
  styleUrls: ['./blog-details.component.css']
})
export class BlogDetailsComponent implements OnInit {
  blog?: Blog;
  blogId!: string;
  comments: Comment[] = [];
  newCommentText = '';

  // tmp jwt
  token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTg3NDExODksInJvbGUiOiJ0b3VyaXN0IiwidXNlcl9pZCI6IjVkOGNmMzcxLTViZmQtNDlkMS04YjI4LWNlMzYzMjE2YTZkYiJ9.lpxuWbl19GUiiX7ucc1bPiZE1UWkM7NkiG_3RlB8J90"
  constructor(
    private route: ActivatedRoute,
    private blogService: BlogService
  ) {}

  ngOnInit(): void {
    this.blogId = this.route.snapshot.paramMap.get('id')!;
    this.loadBlog();
    this.loadComments();
  }

  loadBlog(): void {
    this.blogService.getById(this.blogId).subscribe({
      next: (data) => this.blog = data,
      error: (err) => console.error(err)
    });
  }

  loadComments(): void {
    this.blogService.getComments(this.blogId).subscribe({
      next: (data) => this.comments = data,
      error: (err) => console.error(err)
    });
  }

  likeBlog(): void {
    this.blogService.like(this.blogId, this.token).subscribe({
      next: (likes) => { if (this.blog) this.blog.likes = likes; },
      error: (err) => console.error(err)
    });
  }

  unlikeBlog(): void {
    this.blogService.unlike(this.blogId, this.token).subscribe({
      next: (likes) => { if (this.blog) this.blog.likes = likes; },
      error: (err) => console.error(err)
    });
  }

  addComment(): void {
    if (!this.newCommentText.trim()) return;

    const comment: Comment = {
      blogId: this.blogId,
      text: this.newCommentText.trim()
    };

    this.blogService.createComment(this.blogId, comment, this.token).subscribe({
      next: (c) => {
        this.comments.push(c);
        this.newCommentText = '';
      },
      error: (err) => console.error(err)
    });
  }
}
