import { useEffect, useState } from "react";
import { Button, Layout } from "antd";
import { Route, Routes, useNavigate } from 'react-router-dom';

import { routes } from './Routes.tsx';


import "./App.css";
import "antd/dist/reset.css";

import {
    AppleOutlined,
    CoffeeOutlined,
    OrderedListOutlined, PlusOutlined,
} from "@ant-design/icons";


const { Footer, Content } = Layout;


const App = () => {
  // let safeTop = 0;
  // let safeBottom = 0;


  useEffect(() => {

  }, []);


  const navigate = useNavigate();

  const handleNavigation = (path: string) => {
    navigate(path);
  };


  return (
    <Layout
      style={{
        height: "100vh",
        display: "flex",
        flexDirection: "column",
      }}
    >
      {/* Контент страницы */}
      <Content
        style={{
          flex: 1,
          backgroundColor: "#f0f2f5",
          overflow: "hidden", // Предотвращение прокрутки контента
          display: "flex",
          flexDirection: "column",
          width: "100%",
        }}
      >
        
        {/* Роутер отображает экраны внутри меню */}
        <Routes>
          {routes.map(({ path, Component }: { path: string; Component: React.ComponentType }) => (
            <Route key={path} path={path} element={<Component />} />
          ))}
        </Routes>


      </Content>

      {/* Меню внизу */}
      <Footer
        style={{
          backgroundColor: "#fff",
          display: "flex",
          justifyContent: "space-around", // padding: '10px 0',
          borderTop: "1px solid #e8e8e8",
          position: "sticky", // Закрепляем внизу экрана
          bottom: 0,
          padding: 0,
          width: "100%",
          overflow: "hidden", // Исключение скролла в меню
        }}
      >
        <Button
          type="text"
          icon={<PlusOutlined />}
          size="large"
          onClick={() => handleNavigation("/add-check-screen")}
        />

        <Button
          type="text"
          icon={<OrderedListOutlined />}
          size="large"
          onClick={() => handleNavigation("/view-checks-screen")}
        />
      </Footer>
    </Layout>
  );
};

export default App;


// function setActivePage(path: string) {
//   throw new Error("Function not implemented.");
// }

