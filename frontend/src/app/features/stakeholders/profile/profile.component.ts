import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { StakeholdersService } from '../../../core/services/stakeholders.service';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './profile.component.html'
})
export class ProfileComponent implements OnInit {
  user: any = {}; // moze se definisati interface za tipizaciju

  constructor(private stakeholdersService: StakeholdersService) {}

  ngOnInit(): void {
    console.log('Token:', localStorage.getItem('access_token'));

    this.stakeholdersService.getProfile().subscribe({
      next: (data) => {
        this.user = data;
      },
      error: (err) => {
        console.error('Failed to load profile', err);
      }
    });
  }
}