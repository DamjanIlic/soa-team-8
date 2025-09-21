import { Component } from '@angular/core';
import { NgFor, NgClass } from '@angular/common';

interface User {
  id: number;
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  role: string;
  status: string;
  avatarUrl?: string;
}

@Component({
  selector: 'app-user-list',
  standalone: true,
  imports: [NgFor, NgClass],
  templateUrl: './user-list.component.html',
})
export class UserListComponent {
  users: User[] = [
    {
      id: 1,
      firstName: 'Admin',
      lastName: 'User',
      username: 'admin',
      email: 'admin@example.com',
      role: 'Admin',
      status: 'Active',
      avatarUrl: 'https://i.pravatar.cc/40?img=1',
    },
    {
      id: 2,
      firstName: 'John',
      lastName: 'Doe',
      username: 'john_doe',
      email: 'john@example.com',
      role: 'User',
      status: 'Active',
      avatarUrl: 'https://i.pravatar.cc/40?img=2',
    },
    {
      id: 3,
      firstName: 'Jane',
      lastName: 'Doe',
      username: 'jane_doe',
      email: 'jane@example.com',
      role: 'User',
      status: 'Suspended',
      avatarUrl: 'https://i.pravatar.cc/40?img=3',
    },
  ];
}
