import { ApplicationConfig, importProvidersFrom, provideZoneChangeDetection } from '@angular/core';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { HttpClientModule } from '@angular/common/http';

import { provideMarkdown } from 'ngx-markdown';
import { MarkdownModule } from 'ngx-markdown';

// export const appConfig: ApplicationConfig = {
//   providers: [provideZoneChangeDetection({ eventCoalescing: true }), provideRouter(routes), importProvidersFrom(HttpClientModule),
//     provideMarkdown()
//   ]
// };
import { ReactiveFormsModule } from '@angular/forms';
import { provideHttpClient } from '@angular/common/http';

export const appConfig: ApplicationConfig = {
  providers: [provideZoneChangeDetection({ eventCoalescing: true }), 
    provideRouter(routes),
    importProvidersFrom(ReactiveFormsModule),
    provideHttpClient(),
    importProvidersFrom(HttpClientModule),
    MarkdownModule,
    provideMarkdown()
  ]
};
