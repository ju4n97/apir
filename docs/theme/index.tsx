import { Layout as BasicLayout } from '@rspress/core/theme-original';

const Layout = () => (
  <BasicLayout
    afterNavMenu={
      <div className="esquema-nav-actions">    
        <a
          href="/esquema/docs/installation"
          className="esquema-button esquema-button--primary"
        >
          Install
        </a>
      </div>
    }
  />
);

export * from '@rspress/core/theme-original';
export { Layout };
