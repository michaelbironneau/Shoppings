import { NgModule } from '@angular/core';
import { PreloadAllModules, RouterModule, Routes } from '@angular/router';
import { AuthGuardService as AuthGuard } from './shared/services/auth-guard.service';

const routes: Routes = [
  {
    path: 'home',
    canActivate: [AuthGuard],
    loadChildren: () =>
      import('./home/home.module').then((m) => m.HomePageModule),
  },
  {
    path: '',
    redirectTo: 'home',
    pathMatch: 'full',
  },
  {
    path: 'login',
    loadChildren: () =>
      import('./login/login.module').then((m) => m.LoginPageModule),
  },
  {
    path: 'list/:id',
    canActivate: [AuthGuard],
    loadChildren: () =>
      import('./list/list.module').then((m) => m.ListPageModule),
  },
  {
    path: 'shop/:id',
    canActivate: [AuthGuard],
    loadChildren: () =>
      import('./shop/shop.module').then((m) => m.ShopPageModule),
  },
  {
    path: 'scan',
    loadChildren: () =>
      import('./scan/scan.module').then((m) => m.ScanPageModule),
    canActivate: [AuthGuard],
  },
];

@NgModule({
  imports: [
    RouterModule.forRoot(routes, { preloadingStrategy: PreloadAllModules }),
  ],
  exports: [RouterModule],
})
export class AppRoutingModule {}
