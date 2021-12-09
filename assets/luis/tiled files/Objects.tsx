<?xml version="1.0" encoding="UTF-8"?>
<tileset version="1.5" tiledversion="1.7.2" name="Objects" tilewidth="112" tileheight="128" tilecount="33" columns="0">
 <grid orientation="orthogonal" width="1" height="1"/>
 <tile id="0">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="48" height="64" source="../computers/controlStation1.png"/>
 </tile>
 <tile id="1">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="112" height="64" source="../computers/controlStation2.png"/>
 </tile>
 <tile id="2">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="32" height="64" source="../gas_machines/algaeTank2.png"/>
 </tile>
 <tile id="3">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="48" height="48" source="../gas_machines/gasPump_3x3.png"/>
 </tile>
 <tile id="4">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="../gas_machines/vent1.png"/>
 </tile>
 <tile id="5">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="48" height="64" source="../power_machines/battery3.png"/>
 </tile>
 <tile id="6">
  <properties>
   <property name="needsground" type="bool" value="true"/>
   <property name="wattage" type="int" value="50"/>
  </properties>
  <image width="32" height="48" source="../power_machines/battery4.png"/>
 </tile>
 <tile id="7">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="32" height="64" source="../power_machines/solarGenerator1.png"/>
 </tile>
 <tile id="8">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="64" height="64" source="../power_machines/solarPanel4.png"/>
 </tile>
 <tile id="9">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="32" height="16" source="../plants/plant_1.png"/>
 </tile>
 <tile id="10">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="32" height="32" source="../plants/plant_2.png"/>
 </tile>
 <tile id="11">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="32" height="16" source="../plants/pot_1.png"/>
 </tile>
 <tile id="12">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="80" height="32" source="../furniture/seats/bed1.png"/>
 </tile>
 <tile id="13">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="32" height="32" source="../furniture/seats/chair1.png"/>
 </tile>
 <tile id="14">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="32" height="48" source="../furniture/seats/chair3.png"/>
 </tile>
 <tile id="15">
  <properties>
   <property name="needsroof" type="bool" value="true"/>
  </properties>
  <image width="32" height="16" source="../furniture/decorations/roofScreen.png"/>
 </tile>
 <tile id="16">
  <properties>
   <property name="cxtile" value="decorations/spare-suit"/>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="32" height="48" source="../furniture/decorations/spareSuit.png"/>
 </tile>
 <tile id="17">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="48" height="48" source="../liquid_machines/waterPump_3x3.png"/>
 </tile>
 <tile id="18">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="64" height="80" source="../radar/radar.png"/>
 </tile>
 <tile id="19">
  <properties>
   <property name="needsroof" type="bool" value="true"/>
  </properties>
  <image width="32" height="16" source="../turrets/roofLaser.png"/>
 </tile>
 <tile id="20">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="48" height="32" source="../turrets/Turret1.png"/>
 </tile>
 <tile id="21">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="48" height="32" source="../containers/animalCage.png"/>
 </tile>
 <tile id="22">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="32" height="32" source="../containers/chest.png"/>
 </tile>
 <tile id="23">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="32" height="32" source="../containers/metalBox_blue.png"/>
 </tile>
 <tile id="24">
  <properties>
   <property name="needsground" type="bool" value="true"/>
  </properties>
  <image width="32" height="32" source="../containers/metalBox_orange.png"/>
 </tile>
 <tile id="25">
  <properties>
   <property name="needsroof" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="../furniture/decorations/surveillanceCamera.png"/>
 </tile>
 <tile id="26">
  <properties>
   <property name="cxtile" value="light1"/>
   <property name="powered" type="bool" value="false"/>
   <property name="wattage" type="int" value="-5"/>
  </properties>
  <image width="16" height="16" source="../lights/light_1_off.png"/>
 </tile>
 <tile id="27">
  <properties>
   <property name="cxtile" value="light1"/>
   <property name="powered" type="bool" value="true"/>
   <property name="wattage" type="int" value="-5"/>
  </properties>
  <image width="16" height="16" source="../lights/light_1_on.png"/>
 </tile>
 <tile id="28">
  <properties>
   <property name="needsroof" type="bool" value="true"/>
  </properties>
  <image width="48" height="16" source="../lights/light_2_off.png"/>
 </tile>
 <tile id="29">
  <properties>
   <property name="needsroof" type="bool" value="true"/>
  </properties>
  <image width="48" height="16" source="../lights/light_2_on.png"/>
 </tile>
 <tile id="30">
  <image width="16" height="16" source="../lights/light_3_cable.png"/>
 </tile>
 <tile id="31">
  <properties>
   <property name="needsroof" type="bool" value="true"/>
  </properties>
  <image width="16" height="32" source="../lights/light_3_on.png"/>
 </tile>
 <tile id="32">
  <image width="112" height="128" source="../doors/door_rollingSteel_1.png"/>
 </tile>
</tileset>
