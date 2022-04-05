import { motion, PanInfo } from "framer-motion";
import React, {
  cloneElement,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
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

type ItemAlignerProps = {
  length: number;
  i: number;
};

export interface CanvasProps extends React.PropsWithChildren<any> {
  itemAligner: (props: ItemAlignerProps) => ItemState;
}

export const Canvas = ({
  itemAligner,
  children,
}: React.PropsWithChildren<any>) => {
  const svgRef = useRef<SVGSVGElement>(null);
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
    svgRef.current?.getBoundingClientRect();
    const { width, height } = svgRef.current?.getBoundingClientRect() || {
      width: 0,
      height: 0,
    };
    const dep = width > height ? height : width;
    let newItems = arrayChildren.map((_, i) => ({ x: 0, y: 0 }));
    if (itemAligner) {
      newItems = newItems
        .map((_, i) => itemAligner({ i: i, length: newItems.length }))
        .map(({ x, y }) => ({
          x: x * 0.25 * dep + 0.4 * width,
          y: y * 0.25 * dep + 0.4 * height,
        }));
    }
    setItemState(newItems);
  }, [arrayChildren, svgRef, itemAligner]);

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
        ref={svgRef}
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

export const polygonItemAligner = ({
  length,
  i,
}: ItemAlignerProps): ItemState => {
  console.log(length, i);
  return {
    x: Math.cos((2 * Math.PI * i) / length),
    y: Math.sin((2 * Math.PI * i) / length),
  };
};
