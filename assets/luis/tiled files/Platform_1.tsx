<?xml version="1.0" encoding="UTF-8"?>
<tileset version="1.5" tiledversion="1.7.2" name="Platform_1" tilewidth="16" tileheight="16" tilecount="18" columns="6">
 <image source="../platforms/Platform_1.png" width="96" height="48"/>
 <tile id="0">
  <properties>
   <property name="cxtile" value="platform"/>
  </properties>
 </tile>
 <tile id="1">
  <properties>
   <property name="cxtile" value="platform"/>
  </properties>
 </tile>
 <tile id="2">
  <properties>
   <property name="cxtile" value="platform"/>
  </properties>
 </tile>
 <tile id="3">
  <properties>
   <property name="cxtile" value="smoothplatform"/>
  </properties>
 </tile>
 <tile id="4">
  <properties>
   <property name="cxtile" value="smoothplatform"/>
  </properties>
 </tile>
 <tile id="5">
  <properties>
   <property name="cxtile" value="smoothplatform"/>
  </properties>
 </tile>
 <tile id="6">
  <properties>
   <property name="cxtile" value="platform"/>
  </properties>
 </tile>
 <tile id="7">
  <properties>
   <property name="cxtile" value="platform"/>
  </properties>
 </tile>
 <tile id="8">
  <properties>
   <property name="cxtile" value="platform"/>
  </properties>
 </tile>
 <tile id="9">
  <properties>
   <property name="cxtile" value="smoothplatform"/>
  </properties>
 </tile>
 <tile id="10">
  <properties>
   <property name="cxtile" value="smoothplatform"/>
  </properties>
 </tile>
 <tile id="11">
  <properties>
   <property name="cxtile" value="smoothplatform"/>
  </properties>
 </tile>
 <tile id="12">
  <properties>
   <property name="cxtile" value="platform"/>
  </properties>
 </tile>
 <tile id="13">
  <properties>
   <property name="cxtile" value="platform"/>
  </properties>
 </tile>
 <tile id="14">
  <properties>
   <property name="cxtile" value="platform"/>
  </properties>
 </tile>
 <tile id="15">
  <properties>
   <property name="cxtile" value="smoothplatform"/>
  </properties>
 </tile>
 <tile id="16">
  <properties>
   <property name="cxtile" value="smoothplatform"/>
  </properties>
 </tile>
 <tile id="17">
  <properties>
   <property name="cxtile" value="smoothplatform"/>
  </properties>
 </tile>
 <wangsets>
  <wangset name="Platform" type="edge" tile="6">
   <wangcolor name="Metal platform" color="#ff0000" tile="6" probability="1">
    <properties>
     <property name="ey" type="int" value="0"/>
    </properties>
   </wangcolor>
   <wangtile tileid="1" wangid="0,0,1,0,0,0,1,0"/>
   <wangtile tileid="7" wangid="0,0,1,0,0,0,0,0"/>
   <properties>
    <property name="durability" type="int" value="6"/>
   </properties>
  </wangset>
 </wangsets>
</tileset>
