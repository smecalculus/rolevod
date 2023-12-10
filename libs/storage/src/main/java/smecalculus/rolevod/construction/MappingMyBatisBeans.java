package smecalculus.rolevod.construction;

import static smecalculus.rolevod.configuration.StorageDm.MappingMode.MY_BATIS;

import javax.sql.DataSource;
import org.apache.ibatis.session.SqlSessionFactory;
import org.mybatis.spring.SqlSessionFactoryBean;
import org.mybatis.spring.annotation.MapperScan;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import smecalculus.rolevod.storage.mybatis.UuidTypeHandler;

@ConditionalOnStorageMappingMode(MY_BATIS)
@MapperScan(basePackages = "smecalculus.rolevod.storage.mybatis")
@Configuration(proxyBeanMethods = false)
public class MappingMyBatisBeans {

    @Bean
    public SqlSessionFactory sqlSessionFactory(DataSource dataSource) throws Exception {
        var factoryBean = new SqlSessionFactoryBean();
        factoryBean.setDataSource(dataSource);
        factoryBean.addTypeHandlers(new UuidTypeHandler());
        return factoryBean.getObject();
    }
}
