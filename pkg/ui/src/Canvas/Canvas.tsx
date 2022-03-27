import { motion, PanInfo } from "framer-motion";
import React, { cloneElement, useEffect, useMemo, useState } from "react";
import { Box, Button, ButtonGroup } from "@mui/material";
import {
  AddOutlined,
  FitScreenOutlined,
  RemoveOutlined,
} from "@mui/icons-material";

export interface CanvasItemProps {
  x: number;
  y: number;
  transform: string;
}

type ItemState = {
  x: number;
  y: number;
};

export const Canvas = ({ children }: React.PropsWithChildren<any>) => {
  let arrayChildren = useMemo(
    () => React.Children.toArray(children),
    [children]
  );

  const [itemState, setItemState] = useState<ItemState[]>([]);

  const [scale, setScale] = useState(1);
  const [pan, setPan] = useState({ x: 0, y: 0 });

  const onScale = (e: { deltaY: number }) => {
    setScale(scale + e.deltaY / 500);
    setPan({
      x: pan.x - e.deltaY,
      y: pan.y - e.deltaY,
    });
  };

  const onPan = (info: PanInfo) => {
    setPan({
      x: pan.x + info.delta.x,
      y: pan.y + info.delta.y,
    });
  };

  const zeroView = () => {
    setPan({ x: 0, y: 0 });
    setScale(1);
  };

  useEffect(() => {
    console.log("useEffect");
    if (arrayChildren.length < itemState.length) {
      setItemState(itemState.slice(arrayChildren.length));
    } else {
      const newItems = arrayChildren
        .slice(itemState.length)
        .map((_, i) => ({ x: i * 400, y: i * 400 }));
      setItemState((prevState) => [...prevState, ...newItems]);
    }
  }, [arrayChildren]);

  return (
    <Box
      sx={{
        height: "100%",
        width: "100%",
        position: "relative",
        backgroundColor: "action.disabledBackground",
      }}
    >
      <ButtonGroup
        variant="outlined"
        size="small"
        sx={{ position: "absolute", right: 15, top: 15 }}
        color="secondary"
      >
        <Button onClick={() => setScale((prevState) => prevState + 0.1)}>
          <AddOutlined />
        </Button>
        <Button onClick={() => setScale((prevState) => prevState - 0.1)}>
          <RemoveOutlined />
        </Button>
        <Button onClick={zeroView}>
          <FitScreenOutlined />
        </Button>
      </ButtonGroup>
      <motion.svg
        height="100%"
        width="100%"
        onWheel={(e) => {
          onScale(e);
        }}
        onPan={(e, info) => {
          onPan(info);
        }}
      >
        {itemState.map(({ x, y }, i) => {
          return cloneElement(arrayChildren[i] as React.ReactElement, {
            x: x * scale + pan.x,
            y: y * scale + pan.y,
            transform: `scale(${scale} ${scale})`,
          });
        })}
      </motion.svg>
    </Box>
  );
};
